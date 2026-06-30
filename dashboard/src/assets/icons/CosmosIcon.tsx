import Image from "next/image";
import * as React from "react";
import { memo } from "react";
import CosmosLogo from "@/assets/netbird.svg";

type Props = {
  size?: number;
  className?: string;
};
function CosmosIcon({ size = 16, className }: Props) {
  return (
    <Image
      src={CosmosLogo}
      alt={"Netbird Icon"}
      width={size}
      className={className}
    />
  );
}

export default memo(CosmosIcon);
