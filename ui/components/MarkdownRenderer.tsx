/**
 * File: teamknowl/ui/components/MarkdownRenderer.tsx
 * Purpose: Advanced markdown renderer with support for Wikilinks and Mermaid diagrams.
 * Product/business importance: Provides the core "connected brain" experience for users.
 * 
 * Copyright (c) 2026 John K Johansen
 * License: MIT
 */

"use client";

import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import wikiLink from 'remark-wiki-link';
import Mermaid from './Mermaid';

interface MarkdownRendererProps {
  content: string;
  onLinkClick?: (link: string) => void;
}

const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({ content, onLinkClick }) => {
  return (
    <ReactMarkdown
      remarkPlugins={[
        remarkGfm,
        [wikiLink, { 
          aliasDivider: '|',
          pageResolver: (name: string) => [name.replace(/ /g, '_').toLowerCase()],
          hrefTemplate: (permalink: string) => `/${permalink}`
        }]
      ]}
      components={{
        code({ node, inline, className, children, ...props }: any) {
          const match = /language-mermaid/.exec(className || '');
          if (!inline && match) {
            return (
              <Mermaid chart={String(children).replace(/\n$/, '')} />
            );
          }
          return (
            <code className={className} {...props}>
              {children}
            </code>
          );
        },
        a({ node, href, children, ...props }: any) {
          // Handle internal wiki links
          if (href?.startsWith('/')) {
            return (
              <a
                href={href}
                className="text-indigo-400 hover:text-indigo-300 underline decoration-indigo-500/30 underline-offset-4 cursor-pointer"
                onClick={(e) => {
                  e.preventDefault();
                  if (onLinkClick) {
                    onLinkClick(href.substring(1));
                  }
                }}
                {...props}
              >
                {children}
              </a>
            );
          }
          return (
            <a 
              className="text-indigo-400 hover:text-indigo-300 underline" 
              target="_blank" 
              rel="noopener noreferrer" 
              href={href} 
              {...props}
            >
              {children}
            </a>
          );
        }
      }}
    >
      {content}
    </ReactMarkdown>
  );
};

export default MarkdownRenderer;
