"use client";

import { Users } from "lucide-react";
import PageContainer from "@/layouts/PageContainer";

export default function UsersPage() {
  return (
    <PageContainer>
      <div className="p-6">
        <div className="flex items-center gap-3 mb-6">
          <Users size={28} className="text-kyle shrink-0" />
          <div>
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Users</h1>
            <p className="text-sm text-gray-500 dark:text-nb-gray-400 mt-1">
              Manage users, service accounts, and access roles
            </p>
          </div>
        </div>
        <div className="bg-nb-gray-930/50 rounded-lg border border-nb-gray-900 p-8 text-center">
          <p className="text-nb-gray-400">User management will be available once the backend is connected.</p>
        </div>
      </div>
    </PageContainer>
  );
}
