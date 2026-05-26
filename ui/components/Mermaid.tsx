/**
 * File: teamknowl/ui/components/Mermaid.tsx
 * Purpose: Client-side component for rendering Mermaid.js diagrams.
 * Product/business importance: Enables visual architecture and process mapping within documentation.
 * 
 * Copyright (c) 2026 John K Johansen
 * License: MIT
 */

"use client";

import React, { useEffect, useRef } from 'react';
import mermaid from 'mermaid';

interface MermaidProps {
  chart: string;
}

const Mermaid: React.FC<MermaidProps> = ({ chart }) => {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    mermaid.initialize({
      startOnLoad: true,
      theme: 'dark',
      securityLevel: 'loose',
      fontFamily: 'inherit',
    });

    if (ref.current) {
      mermaid.contentLoaded();
    }
  }, [chart]);

  useEffect(() => {
    if (ref.current) {
      ref.current.removeAttribute('data-processed');
      mermaid.render('mermaid-svg-' + Math.random().toString(36).substr(2, 9), chart).then(
        ({ svg }) => {
          if (ref.current) {
            ref.current.innerHTML = svg;
          }
        }
      ).catch((err) => {
        console.error('Mermaid render error:', err);
      });
    }
  }, [chart]);

  return (
    <div className="mermaid flex justify-center my-8 bg-slate-900/50 p-4 rounded-lg border border-slate-800" ref={ref}>
      {chart}
    </div>
  );
};

export default Mermaid;
