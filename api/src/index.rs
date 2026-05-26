//! File: teamknowl/api/src/index.rs
//! Purpose: High-performance search indexing engine using Tantivy.
//! Product/business importance: enables rapid search and context retrieval across large knowledge bases.
//!
//! Copyright (c) 2026 John K Johansen
//! License: MIT

use anyhow::Result;
use std::path::Path;
use tantivy::schema::*;
use tantivy::{Index, IndexWriter, ReloadPolicy, collector::TopDocs};

pub struct SearchIndex {
    index: Index,
    schema: Schema,
}

impl SearchIndex {
    pub fn new(index_path: &Path) -> Result<Self> {
        let mut schema_builder = Schema::builder();
        schema_builder.add_text_field("title", TEXT | STORED);
        schema_builder.add_text_field("body", TEXT | STORED);
        schema_builder.add_text_field("path", STRING | STORED);
        let schema = schema_builder.build();

        std::fs::create_dir_all(index_path)?;
        let index = Index::open_or_create(
            tantivy::directory::MmapDirectory::open(index_path)?,
            schema.clone(),
        )?;

        Ok(Self { index, schema })
    }

    pub fn get_writer(&self) -> Result<IndexWriter> {
        Ok(self.index.writer(50_000_000)?)
    }

    pub fn search(&self, query_str: &str) -> Result<Vec<String>> {
        let reader = self
            .index
            .reader_builder()
            .reload_policy(ReloadPolicy::OnCommitWithDelay)
            .try_into()?;
        let searcher = reader.searcher();

        let title = self.schema.get_field("title").unwrap();
        let body = self.schema.get_field("body").unwrap();
        let query_parser = tantivy::query::QueryParser::for_index(&self.index, vec![title, body]);
        let query = query_parser.parse_query(query_str)?;

        let top_docs = searcher.search(&query, &TopDocs::with_limit(10).order_by_score())?;

        let mut results = Vec::new();
        for (_score, doc_address) in top_docs {
            let retrieved_doc: tantivy::TantivyDocument = searcher.doc(doc_address)?;
            let path_field = self.schema.get_field("path").unwrap();
            if let Some(path_val) = retrieved_doc.get_first(path_field).and_then(|v| v.as_str()) {
                results.push(path_val.to_string());
            }
        }

        Ok(results)
    }
}
