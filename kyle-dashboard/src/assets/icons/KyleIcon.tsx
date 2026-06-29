import Image from "next/image";
import * as React from "react";
import { memo } from "react";
import KyleLogoMark from "@/assets/kyle.svg";

type Props = {
  size?: number;
  className?: string;
};
function KyleIcon({ size = 16, className }: Props) {
  return (
    <Image
      src={KyleLogoMark}
      alt={"Kyle Icon"}
      width={size}
      className={className}
    />
  );
}

export default memo(KyleIcon);
