"use client";

import { Terminal } from "lucide-react";
import PageContainer from "@/layouts/PageContainer";

export default function ConnectionsPage() {
  return (
    <PageContainer>
      <div className="p-6">
        <div className="flex items-center gap-3 mb-6">
          <Terminal size={28} className="text-kyle shrink-0" />
          <div>
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Connections</h1>
            <p className="text-sm text-gray-500 dark:text-nb-gray-400 mt-1">
              Active remote desktop and SSH sessions
            </p>
          </div>
        </div>
        <div className="bg-nb-gray-930/50 rounded-lg border border-nb-gray-900 p-8 text-center">
          <p className="text-nb-gray-400">No active connections. Connect to a server to start a session.</p>
        </div>
      </div>
    </PageContainer>
  );
}
