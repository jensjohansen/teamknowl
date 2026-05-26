//! File: teamknowl/api/src/storage.rs
//! Purpose: S3-compatible storage client for note retrieval.
//! Product/business importance: enables durable, scalable storage of knowledge artifacts.
//!
//! Copyright (c) 2026 John K Johansen
//! License: MIT

use anyhow::Result;
use aws_config::meta::region::RegionProviderChain;
use aws_sdk_s3::Client;
use lru::LruCache;
use std::num::NonZeroUsize;
use tokio::sync::Mutex;

pub struct StorageClient {
    client: Client,
    bucket: String,
    cache: Mutex<LruCache<String, String>>,
}

impl StorageClient {
    pub async fn new(bucket: String) -> Result<Self> {
        let region_provider = RegionProviderChain::default_provider().or_else("us-east-1");
        let mut config_loader = aws_config::defaults(aws_config::BehaviorVersion::latest())
            .region(region_provider);

        // Check for custom S3 endpoint (for Ceph/Minio)
        if let Ok(endpoint) = std::env::var("S3_ENDPOINT") {
            config_loader = config_loader.endpoint_url(endpoint);
        }

        // Support standard AWS env vars if provided via Operator (S3_ACCESS_KEY -> AWS_ACCESS_KEY_ID)
        if let Ok(access_key) = std::env::var("S3_ACCESS_KEY") {
            std::env::set_var("AWS_ACCESS_KEY_ID", access_key);
        }
        if let Ok(secret_key) = std::env::var("S3_SECRET_KEY") {
            std::env::set_var("AWS_SECRET_ACCESS_KEY", secret_key);
        }

        let config = config_loader.load().await;
        let client = Client::new(&config);

        // Initialize LRU cache with a capacity of 1000 items
        let cache = Mutex::new(LruCache::new(NonZeroUsize::new(1000).unwrap()));

        Ok(Self {
            client,
            bucket,
            cache,
        })
    }

    pub async fn get_object(&self, key: &str) -> Result<String> {
        // 1. Check cache
        {
            let mut cache = self.cache.lock().await;
            if let Some(content) = cache.get(key) {
                return Ok(content.clone());
            }
        }

        // 2. Fetch from S3
        let resp = self
            .client
            .get_object()
            .bucket(&self.bucket)
            .key(key)
            .send()
            .await?;

        let data = resp.body.collect().await?.into_bytes();
        let content = String::from_utf8(data.to_vec())?;

        // 3. Update cache
        {
            let mut cache = self.cache.lock().await;
            cache.put(key.to_string(), content.clone());
        }

        Ok(content)
    }
}
