import Image from "next/image";
import * as React from "react";
import { memo } from "react";
import CosmosLogoMark from "@/assets/cosmos.svg";

type Props = {
  size?: number;
  className?: string;
};
function CosmosIcon({ size = 16, className }: Props) {
  return (
    <Image
      src={CosmosLogoMark}
      alt={"Cosmos Icon"}
      width={size}
      className={className}
    />
  );
}

export default memo(CosmosIcon);
