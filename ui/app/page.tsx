/**
 * File: teamknowl/ui/app/page.tsx
 * Purpose: Main entry point for the TeamKnowl Obsidian-like interface.
 * Product/business importance: Provides the primary user interface for knowledge discovery and navigation.
 * 
 * Copyright (c) 2026 John K Johansen
 * License: MIT
 */

"use client";

import { useEffect, useState, useMemo, useCallback } from "react";
import { 
  Settings, 
  FolderClosed, 
  ChevronRight, 
  Sparkles, 
  FileText, 
  Loader2,
  Network,
  BookOpen,
  Search
} from "lucide-react";
import MarkdownRenderer from "../components/MarkdownRenderer";
import GraphView from "../components/GraphView";

interface NoteInfo {
  id: string;
  name: string;
}

export default function Home() {
  const [notes, setNotes] = useState<NoteInfo[]>([]);
  const [selectedNoteId, setSelectedNoteId] = useState<string | null>(null);
  const [content, setContent] = useState<string>("");
  const [viewMode, setViewMode] = useState<'document' | 'graph'>('document');
  const [loading, setLoading] = useState(true);
  const [contentLoading, setContentLoading] = useState(false);

  // Mock graph data for visualization demonstration
  const graphData = useMemo(() => {
    const nodes = notes.map(n => ({ id: n.id, name: n.name, val: 1 }));
    const links = notes.length > 1 ? [
      { source: notes[0].id, target: notes[1].id },
      ...(notes.length > 2 ? [{ source: notes[1].id, target: notes[2].id }] : [])
    ] : [];
    return { nodes, links };
  }, [notes]);

  useEffect(() => {
    const fetchNotes = async () => {
      try {
        // In a real scenario, this would call the Rust API search/list endpoint
        // For now, we use the existing /api/v1/list or mock if unavailable
        const response = await fetch("/api/v1/list");
        if (response.ok) {
          const data = await response.json();
          const formatted = data.map((f: any) => ({
            id: f.path.replace('.md', ''),
            name: f.name.replace('.md', '')
          }));
          setNotes(formatted);
          if (formatted.length > 0) {
            handleNoteSelect(formatted[0].id);
          }
        } else {
          // Mock data for development if API is not running
          const mockNotes = [
            { id: 'welcome', name: 'Welcome to TeamKnowl' },
            { id: 'architecture', name: 'System Architecture' },
            { id: 'deployment', name: 'Deployment Guide' }
          ];
          setNotes(mockNotes);
          handleNoteSelect(mockNotes[0].id);
        }
      } catch (error) {
        console.error("Failed to fetch notes:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchNotes();
  }, [handleNoteSelect]);

  const handleNoteSelect = useCallback(async (noteId: string) => {
    setSelectedNoteId(noteId);
    setContentLoading(true);
    try {
      // Calling the new Rust API endpoint (proxied via Next.js or direct)
      const response = await fetch(`/api/v1/context/${noteId}`);
      if (response.ok) {
        const text = await response.text();
        setContent(text);
      } else {
        setContent(`# ${noteId}\n\nThis is a placeholder for the note content. The Rust API endpoint \`/v1/context/${noteId}\` returned ${response.status}.`);
      }
    } catch (error) {
      console.error("Failed to fetch content:", error);
      setContent("# Error\nFailed to load content from the knowledge engine.");
    } finally {
      setContentLoading(false);
      if (viewMode === 'graph') setViewMode('document');
    }
  }, [viewMode]);

  const selectedNote = notes.find(n => n.id === selectedNoteId);

  return (
    <main className="flex h-screen overflow-hidden bg-slate-950 text-slate-200">
      {/* Sidebar */}
      <aside className="w-64 border-r border-slate-800 bg-slate-900/50 flex flex-col">
        <div className="p-4 flex items-center gap-2 border-b border-slate-800">
          <div className="w-8 h-8 bg-indigo-600 rounded flex items-center justify-center font-bold text-white">K</div>
          <h1 className="font-bold text-slate-100 tracking-tight text-lg">TeamKnowl</h1>
        </div>

        <div className="p-4">
          <div className="relative">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-slate-500" />
            <input 
              type="text" 
              placeholder="Search notes..." 
              className="w-full bg-slate-950 border border-slate-800 rounded-md py-2 pl-9 pr-4 text-sm focus:outline-none focus:ring-1 focus:ring-indigo-500 transition-all"
            />
          </div>
        </div>

        <nav className="flex-1 p-2 space-y-1 overflow-y-auto">
          <div className="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
            Knowledge Base
          </div>
          
          {loading ? (
            <div className="flex items-center justify-center p-4">
              <Loader2 className="w-5 h-5 animate-spin text-slate-500" />
            </div>
          ) : (
            notes.map((note) => (
              <div
                key={note.id}
                onClick={() => handleNoteSelect(note.id)}
                className={`
                  flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors cursor-pointer
                  ${selectedNoteId === note.id ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/20' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200'}
                `}
              >
                <FileText size={16} />
                <span className="truncate">{note.name}</span>
              </div>
            ))
          )}
        </nav>

        <div className="p-4 border-t border-slate-800 space-y-2">
          <button 
            onClick={() => setViewMode('graph')}
            className={`w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors ${viewMode === 'graph' ? 'bg-indigo-600 text-white' : 'text-slate-400 hover:bg-slate-800/50'}`}
          >
            <Network size={18} />
            <span>Graph View</span>
          </button>
          <div className="flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium text-slate-400 hover:bg-slate-800/50 cursor-pointer">
            <Settings size={18} />
            <span>Settings</span>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col min-w-0">
        <header className="h-12 border-b border-slate-800 flex items-center justify-between px-4 bg-slate-900/20">
          <div className="flex items-center gap-2 text-sm text-slate-400 truncate">
            <FolderClosed size={14} />
            <ChevronRight size={14} />
            <span className="truncate">Notes</span>
            <ChevronRight size={14} />
            <span className="text-slate-200 font-medium truncate">{selectedNote?.name || 'Select a note'}</span>
          </div>
          <div className="flex items-center gap-3">
            {viewMode === 'graph' && (
              <button 
                onClick={() => setViewMode('document')}
                className="flex items-center gap-2 px-3 py-1 rounded-md bg-slate-800 text-slate-300 text-xs font-medium border border-slate-700 hover:bg-slate-700 transition-colors"
              >
                <BookOpen size={14} />
                Close Graph
              </button>
            )}
            <button className="flex items-center gap-2 px-3 py-1 rounded-md bg-indigo-600/10 text-indigo-400 text-xs font-medium border border-indigo-500/20 hover:bg-indigo-600/20 transition-colors">
              <Sparkles size={14} />
              AI Context
            </button>
          </div>
        </header>

        <div className="flex-1 overflow-hidden relative">
          {viewMode === 'graph' ? (
            <div className="absolute inset-0">
              <GraphView 
                data={graphData} 
                onNodeClick={(id) => handleNoteSelect(id)}
              />
            </div>
          ) : (
            <div className="h-full overflow-y-auto p-8 lg:p-12">
              <div className="max-w-4xl mx-auto w-full">
                {contentLoading ? (
                  <div className="flex items-center justify-center h-64">
                    <Loader2 className="w-8 h-8 animate-spin text-indigo-500" />
                  </div>
                ) : (
                  <article className="prose prose-invert prose-indigo max-w-none">
                    <MarkdownRenderer 
                      content={content} 
                      onLinkClick={(id) => handleNoteSelect(id)} 
                    />
                  </article>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </main>
  );
}
