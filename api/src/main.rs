//! File: teamknowl/api/src/main.rs
//! Purpose: High-performance Rust-native API for TeamKnowl knowledge management.
//! Product/business importance: provides the foundation for S3-backed note serving and AI context retrieval.
//!
//! Copyright (c) 2026 John K Johansen
//! License: MIT

use api::index::SearchIndex;
use api::storage::StorageClient;
use api::resolution::ContextResolver;
use axum::{
    extract::{Path as AxumPath, State},
    routing::get,
    Router,
};
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

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize logging
    tracing_subscriber::fmt::init();

    info!("Starting TeamKnowl API...");

    let bucket = std::env::var("TEAMKNOWL_BUCKET").unwrap_or_else(|_| "teamknowl-notes".to_string());
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
        .with_state(state);

    // Run it with hyper on localhost:3000
    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
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
