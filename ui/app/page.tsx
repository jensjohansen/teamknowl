import { Search, Hash, Clock, Settings, FolderClosed, ChevronRight, Share2, Sparkles } from "lucide-react";

export default function Home() {
  return (
    <main className="flex h-screen overflow-hidden">
      {/* Sidebar - Navigation */}
      <aside className="w-64 border-r border-slate-800 bg-slate-900/50 flex flex-col">
        <div className="p-4 flex items-center gap-2 border-b border-slate-800">
          <div className="w-8 h-8 rounded bg-indigo-500 flex items-center justify-center font-bold text-white shadow-lg shadow-indigo-500/20">
            K
          </div>
          <h1 className="font-bold text-slate-100 tracking-tight">TeamKnowl</h1>
        </div>

        <nav className="flex-1 p-2 space-y-1 overflow-y-auto">
          <div className="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
            Navigation
          </div>
          <NavItem icon={<Search size={18} />} label="Search" active />
          <NavItem icon={<Clock size={18} />} label="Recent" />
          <NavItem icon={<FolderClosed size={18} />} label="Documents" />
          
          <div className="px-3 py-2 mt-4 text-xs font-semibold text-slate-500 uppercase tracking-wider">
            Knowledge Bases
          </div>
          <NavItem icon={<Hash size={18} />} label="devops-knowl" />
          <NavItem icon={<Hash size={18} />} label="support-knowl" />
          <NavItem icon={<Hash size={18} />} label="team-knowl" />
        </nav>

        <div className="p-4 border-t border-slate-800">
          <NavItem icon={<Settings size={18} />} label="Settings" />
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col bg-slate-950">
        {/* Toolbar */}
        <header className="h-12 border-b border-slate-800 flex items-center justify-between px-4 bg-slate-900/20">
          <div className="flex items-center gap-2 text-sm text-slate-400">
            <FolderClosed size={14} />
            <ChevronRight size={14} />
            <span>Architecture</span>
            <ChevronRight size={14} />
            <span className="text-slate-200 font-medium">SystemOverview.md</span>
          </div>
          <div className="flex items-center gap-3">
            <button className="flex items-center gap-2 px-3 py-1 rounded-md bg-indigo-600/10 text-indigo-400 text-xs font-medium border border-indigo-500/20 hover:bg-indigo-600/20 transition-colors">
              <Sparkles size={14} />
              AI Insights
            </button>
            <button className="text-slate-400 hover:text-slate-200 transition-colors">
              <Share2 size={18} />
            </button>
          </div>
        </header>

        {/* Editor / Reader */}
        <div className="flex-1 overflow-y-auto p-12 max-w-4xl mx-auto w-full prose prose-invert prose-slate">
          <h1 className="text-4xl font-bold mb-8">System Overview</h1>
          
          <p className="text-xl text-slate-400 leading-relaxed mb-6 italic">
            "An open-source, Kubernetes-native knowledge base platform designed for AI + Human teams."
          </p>

          <h2 className="text-2xl font-semibold mt-12 mb-4">Architecture</h2>
          <p className="mb-4">
            TeamKnowl is built to handle massive scale documentation repositories by decoupling 
            the synchronization logic (git-sync) from the retrieval layer (API) and the visualization 
            layer (UI).
          </p>

          <div className="bg-slate-900 p-6 rounded-lg border border-slate-800 font-mono text-sm mb-6">
            <div className="text-indigo-400 mb-2"># Distributed Search Logic</div>
            <div>[Bleve] Indexing enabled...</div>
            <div>[S3] Connecting to CEPH Object Store...</div>
            <div>[Auth] Validating RBAC for jc01...</div>
          </div>

          <h2 className="text-2xl font-semibold mt-12 mb-4">Connected Knowledge</h2>
          <div className="grid grid-cols-2 gap-4">
            <div className="p-4 border border-slate-800 rounded-lg hover:border-indigo-500/50 transition-colors cursor-pointer group">
              <div className="text-sm font-medium text-slate-300 group-hover:text-indigo-400 mb-1">Backlink</div>
              <div className="text-xs text-slate-500 uppercase tracking-tighter">PRD.md</div>
            </div>
            <div className="p-4 border border-slate-800 rounded-lg hover:border-indigo-500/50 transition-colors cursor-pointer group">
              <div className="text-sm font-medium text-slate-300 group-hover:text-indigo-400 mb-1">Backlink</div>
              <div className="text-xs text-slate-500 uppercase tracking-tighter">ImplementationPlan.md</div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}

function NavItem({ icon, label, active = false }: { icon: React.ReactNode, label: string, active?: boolean }) {
  return (
    <div className={`
      flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors cursor-pointer
      ${active ? 'bg-indigo-600/10 text-indigo-400' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200'}
    `}>
      {icon}
      <span>{label}</span>
    </div>
  );
}
