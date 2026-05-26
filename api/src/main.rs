//! File: teamknowl/api/src/main.rs
//! Purpose: High-performance Rust-native API for TeamKnowl knowledge management.
//! Product/business importance: provides the foundation for S3-backed note serving and AI context retrieval.
//!
//! Copyright (c) 2026 John K Johansen
//! License: MIT

use api::index::SearchIndex;
use api::resolution::ContextResolver;
use api::storage::StorageClient;
use axum::{
    extract::{Path as AxumPath, State},
    routing::{get, post},
    Json, Router,
};
use serde::{Deserialize, Serialize};
use std::net::SocketAddr;
use std::path::Path;
use std::sync::Arc;
use tracing::info;

struct AppState {
    #[allow(dead_code)]
    storage: StorageClient,
    #[allow(dead_code)]
    index: SearchIndex,
    resolver: ContextResolver,
}

#[derive(Deserialize)]
struct BatchContextRequest {
    note_ids: Vec<String>,
}

#[derive(Serialize)]
struct BatchContextResponse {
    results: Vec<NoteContext>,
}

#[derive(Serialize)]
struct NoteContext {
    id: String,
    content: Option<String>,
    error: Option<String>,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize logging
    tracing_subscriber::fmt::init();

    info!("Starting TeamKnowl API...");

    let bucket = std::env::var("S3_BUCKET")
        .or_else(|_| std::env::var("TEAMKNOWL_BUCKET"))
        .unwrap_or_else(|_| "teamknowl-notes".to_string());
    
    let index_path =
        std::env::var("TEAMKNOWL_INDEX_PATH").unwrap_or_else(|_| "data/index".to_string());

    let storage = StorageClient::new(bucket).await?;
    let index = SearchIndex::new(Path::new(&index_path))?;
    let resolver = ContextResolver::new();

    let state = Arc::new(AppState {
        storage,
        index,
        resolver,
    });

    // Build our application with routes
    let app = Router::new()
        .route("/health", get(health_check))
        .route("/v1/context/:note_id", get(get_context))
        .route("/v1/batch/context", post(batch_get_context))
        .with_state(state);

    // Run it with hyper on localhost:8080
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    info!("Listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

async fn health_check() -> &'static str {
    "OK"
}

async fn get_context(
    AxumPath(note_id): AxumPath<String>,
    State(state): State<Arc<AppState>>,
) -> Result<String, String> {
    state
        .resolver
        .resolve_flat_context(&state.storage, &note_id)
        .await
        .map_err(|e| e.to_string())
}

async fn batch_get_context(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<BatchContextRequest>,
) -> Json<BatchContextResponse> {
    let results = state
        .resolver
        .resolve_batch_context(&state.storage, payload.note_ids)
        .await;

    let formatted_results = results
        .into_iter()
        .map(|(id, res)| match res {
            Ok(content) => NoteContext {
                id,
                content: Some(content),
                error: None,
            },
            Err(e) => NoteContext {
                id,
                content: None,
                error: Some(e.to_string()),
            },
        })
        .collect();

    Json(BatchContextResponse {
        results: formatted_results,
    })
}
