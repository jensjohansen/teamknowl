"use client";

import { useEffect, useState } from "react";
import { Search, Hash, Clock, Settings, FolderClosed, ChevronRight, Share2, Sparkles, FileText, Loader2 } from "lucide-react";
import ReactMarkdown from "react-markdown";

interface FileInfo {
  name: string;
  path: string;
}

export default function Home() {
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [selectedFile, setSelectedFile] = useState<FileInfo | null>(null);
  const [content, setContent] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [contentLoading, setContentLoading] = useState(false);

  useEffect(() => {
    fetchFiles();
  }, []);

  const fetchFiles = async () => {
    try {
      const response = await fetch("/api/v1/list");
      const data = await response.json();
      setFiles(data);
      if (data.length > 0) {
        handleFileSelect(data[0]);
      }
    } catch (error) {
      console.error("Failed to fetch files:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleFileSelect = async (file: FileInfo) => {
    setSelectedFile(file);
    setContentLoading(true);
    try {
      const response = await fetch(`/api/v1/context?path=${encodeURIComponent(file.path)}`);
      const text = await response.text();
      setContent(text);
    } catch (error) {
      console.error("Failed to fetch content:", error);
      setContent("# Error\nFailed to load content.");
    } finally {
      setContentLoading(false);
    }
  };

  return (
    <main className="flex h-screen overflow-hidden bg-slate-950 text-slate-200">
      {/* Sidebar */}
      <aside className="w-64 border-r border-slate-800 bg-slate-900/50 flex flex-col">
        <div className="p-4 flex items-center gap-2 border-b border-slate-800">
          <img src="/logo.png" alt="Knowl" className="w-8 h-8 rounded" />
          <h1 className="font-bold text-slate-100 tracking-tight text-lg">TeamKnowl</h1>
        </div>

        <nav className="flex-1 p-2 space-y-1 overflow-y-auto">
          <div className="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
            Documents
          </div>
          
          {loading ? (
            <div className="flex items-center justify-center p-4">
              <Loader2 className="w-5 h-5 animate-spin text-slate-500" />
            </div>
          ) : (
            files.map((file) => (
              <div
                key={file.path}
                onClick={() => handleFileSelect(file)}
                className={`
                  flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors cursor-pointer
                  ${selectedFile?.path === file.path ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/20' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200'}
                `}
              >
                <FileText size={16} />
                <span className="truncate">{file.name}</span>
              </div>
            ))
          )}
        </nav>

        <div className="p-4 border-t border-slate-800">
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
            <span className="truncate">{selectedFile?.path.split('/').slice(0, -1).join(' / ') || 'Root'}</span>
            <ChevronRight size={14} />
            <span className="text-slate-200 font-medium truncate">{selectedFile?.name || 'Select a file'}</span>
          </div>
          <div className="flex items-center gap-3">
            <button className="flex items-center gap-2 px-3 py-1 rounded-md bg-indigo-600/10 text-indigo-400 text-xs font-medium border border-indigo-500/20 hover:bg-indigo-600/20 transition-colors">
              <Sparkles size={14} />
              AI Insights
            </button>
          </div>
        </header>

        <div className="flex-1 overflow-y-auto p-8 lg:p-12">
          <div className="max-w-4xl mx-auto w-full">
            {contentLoading ? (
              <div className="flex items-center justify-center h-64">
                <Loader2 className="w-8 h-8 animate-spin text-indigo-500" />
              </div>
            ) : (
              <article className="prose prose-invert prose-indigo max-w-none">
                <ReactMarkdown>{content}</ReactMarkdown>
              </article>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
