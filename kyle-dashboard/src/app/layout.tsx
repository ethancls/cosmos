import { globalMetaTitle } from "@utils/meta";
import type { Metadata } from "next";
import AppLayout from "@/layouts/AppLayout";

export const metadata: Metadata = {
  title: `${globalMetaTitle}`,
  description:
    "Kyle — Secure remote access gateway. SSH, RDP, VNC in your browser with zero trust bastion architecture.",
};
export default AppLayout;
