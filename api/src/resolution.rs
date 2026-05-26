//! File: teamknowl/api/src/resolution.rs
//! Purpose: Logic for resolving [[Wikilinks]] and flattening note context for LLMs.
//! Product/business importance: enables AI agents to understand the connections between knowledge artifacts.
//!
//! Copyright (c) 2026 John K Johansen
//! License: MIT

use regex::Regex;
use std::collections::HashSet;
use crate::storage::StorageClient;
use anyhow::Result;

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
    /// the content of its directly linked notes.
    pub async fn resolve_flat_context(&self, storage: &StorageClient, note_id: &str) -> Result<String> {
        let mut flat_context = String::new();
        let mut visited = HashSet::new();

        // 1. Fetch the primary note
        let main_content = storage.get_object(note_id).await?;
        flat_context.push_str(&format!("# Primary Note: {}\n\n", note_id));
        flat_context.push_str(&main_content);
        flat_context.push_str("\n\n---\n\n");

        visited.insert(note_id.to_string());

        // 2. Extract links and fetch related content (1-level depth for MVP)
        let links = self.extract_links(&main_content);
        for link in links {
            if !visited.contains(&link) {
                if let Ok(linked_content) = storage.get_object(&link).await {
                    flat_context.push_str(&format!("## Linked Note: {}\n\n", link));
                    flat_context.push_str(&linked_content);
                    flat_context.push_str("\n\n---\n\n");
                    visited.insert(link);
                }
            }
        }

        Ok(flat_context)
    }
}
