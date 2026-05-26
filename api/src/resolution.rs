//! File: teamknowl/api/src/resolution.rs
//! Purpose: Logic for resolving [[Wikilinks]] and flattening note context for LLMs.
//! Product/business importance: enables AI agents to understand the connections between knowledge artifacts.
//!
//! Copyright (c) 2026 John K Johansen
//! License: MIT

use crate::storage::StorageClient;
use anyhow::Result;
use futures::future::join_all;
use regex::Regex;
use std::collections::HashSet;

pub struct ContextResolver {
    link_re: Regex,
}

impl ContextResolver {
    pub fn new() -> Self {
        Self {
            link_re: Regex::new(r"\[\[(.*?)\]\]").unwrap(),
        }
    }

    /// Extracts all wikilinks from a markdown string.
    pub fn extract_links(&self, content: &str) -> Vec<String> {
        self.link_re
            .captures_iter(content)
            .map(|cap| cap[1].to_string())
            .collect()
    }

    /// Generates a "flat" context by taking a starting note and appending
    /// the content of its directly linked notes in parallel.
    pub async fn resolve_flat_context(&self, storage: &StorageClient, note_id: &str) -> Result<String> {
        let mut flat_context = String::new();
        let mut visited = HashSet::new();

        // 1. Fetch the primary note
        let main_content = storage.get_object(note_id).await?;
        flat_context.push_str(&format!("# Primary Note: {}\n\n", note_id));
        flat_context.push_str(&main_content);
        flat_context.push_str("\n\n---\n\n");

        visited.insert(note_id.to_string());

        // 2. Extract links and fetch related content in parallel
        let links: Vec<String> = self
            .extract_links(&main_content)
            .into_iter()
            .filter(|link| !visited.contains(link))
            .collect();

        let fetch_futures = links.iter().map(|link| async move {
            match storage.get_object(link).await {
                Ok(content) => Some((link.clone(), content)),
                Err(_) => None,
            }
        });

        let results = join_all(fetch_futures).await;

        for result in results.into_iter().flatten() {
            let (link, linked_content) = result;
            flat_context.push_str(&format!("## Linked Note: {}\n\n", link));
            flat_context.push_str(&linked_content);
            flat_context.push_str("\n\n---\n\n");
            visited.insert(link);
        }

        Ok(flat_context)
    }

    /// Resolves context for multiple notes in parallel.
    pub async fn resolve_batch_context(
        &self,
        storage: &StorageClient,
        note_ids: Vec<String>,
    ) -> Vec<(String, Result<String>)> {
        let fetch_futures = note_ids.into_iter().map(|id| async move {
            let res = self.resolve_flat_context(storage, &id).await;
            (id, res)
        });

        join_all(fetch_futures).await
    }
}
