//! File: teamknowl/api/src/storage.rs
//! Purpose: S3-compatible storage client for note retrieval.
//! Product/business importance: enables durable, scalable storage of knowledge artifacts.
//!
//! Copyright (c) 2026 John K Johansen
//! License: MIT

use anyhow::Result;
use aws_config::meta::region::RegionProviderChain;
use aws_sdk_s3::Client;

pub struct StorageClient {
    client: Client,
    bucket: String,
}

impl StorageClient {
    pub async fn new(bucket: String) -> Result<Self> {
        let region_provider = RegionProviderChain::default_provider().or_else("us-east-1");
        let config = aws_config::defaults(aws_config::BehaviorVersion::latest())
            .region(region_provider)
            .load()
            .await;
        let client = Client::new(&config);

        Ok(Self { client, bucket })
    }

    pub async fn get_object(&self, key: &str) -> Result<String> {
        let resp = self
            .client
            .get_object()
            .bucket(&self.bucket)
            .key(key)
            .send()
            .await?;

        let data = resp.body.collect().await?.into_bytes();
        Ok(String::from_utf8(data.to_vec())?)
    }
}
