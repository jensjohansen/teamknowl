import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "TeamKnowl | AI + Human Knowledge Base",
  description: "Kubernetes-native, Obsidian-like knowledge base for modern teams.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className="bg-slate-950 text-slate-200 antialiased">
        {children}
      </body>
    </html>
  );
}
