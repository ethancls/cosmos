"use client";

import { ShieldCheck } from "lucide-react";
import PageContainer from "@/layouts/PageContainer";

export default function PoliciesPage() {
  return (
    <PageContainer>
      <div className="p-6">
        <div className="flex items-center gap-3 mb-6">
          <ShieldCheck size={28} className="text-kyle shrink-0" />
          <div>
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Policies</h1>
            <p className="text-sm text-gray-500 dark:text-nb-gray-400 mt-1">
              Zero trust access policies for your infrastructure
            </p>
          </div>
        </div>
        <div className="bg-nb-gray-930/50 rounded-lg border border-nb-gray-900 p-8 text-center">
          <p className="text-nb-gray-400">No policies defined yet. Create your first access policy to enforce zero trust.</p>
        </div>
      </div>
    </PageContainer>
  );
}
