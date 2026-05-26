/**
 * File: teamknowl/ui/components/GraphView.tsx
 * Purpose: Interactive knowledge graph visualization using Force-Directed Graphs.
 * Product/business importance: Provides a visual overview of document relationships and knowledge clusters.
 * 
 * Copyright (c) 2026 John K Johansen
 * License: MIT
 */

"use client";

import React, { useMemo } from 'react';
import ForceGraph2D from 'react-force-graph-2d';

interface Node {
  id: string;
  name: string;
  val: number;
}

interface Link {
  source: string;
  target: string;
}

interface GraphData {
  nodes: Node[];
  links: Link[];
}

interface GraphViewProps {
  data: GraphData;
  onNodeClick?: (nodeId: string) => void;
  width?: number;
  height?: number;
}

const GraphView: React.FC<GraphViewProps> = ({ data, onNodeClick, width, height }) => {
  const processedData = useMemo(() => data, [data]);

  return (
    <div className="bg-slate-950 rounded-lg border border-slate-800 overflow-hidden">
      <ForceGraph2D
        graphData={processedData}
        nodeLabel="name"
        nodeColor={() => '#6366f1'} // indigo-500
        linkColor={() => '#334155'} // slate-700
        backgroundColor="#020617" // slate-950
        width={width}
        height={height}
        onNodeClick={(node: any) => onNodeClick && onNodeClick(node.id)}
        nodeCanvasObject={(node: any, ctx, globalScale) => {
          const label = node.name;
          const fontSize = 12 / globalScale;
          ctx.font = `${fontSize}px Sans-Serif`;
          const textWidth = ctx.measureText(label).width;
          const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2) as [number, number];

          ctx.fillStyle = 'rgba(2, 6, 23, 0.8)';
          ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, bckgDimensions[0], bckgDimensions[1]);

          ctx.textAlign = 'center';
          ctx.textBaseline = 'middle';
          ctx.fillStyle = node.color || '#94a3b8'; // slate-400
          ctx.fillText(label, node.x, node.y);

          node.__bckgDimensions = bckgDimensions; // to use in nodePointerAreaPaint
        }}
        nodePointerAreaPaint={(node: any, color, ctx) => {
          ctx.fillStyle = color;
          const bckgDimensions = node.__bckgDimensions;
          bckgDimensions && ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, bckgDimensions[0], bckgDimensions[1]);
        }}
      />
    </div>
  );
};

export default GraphView;
