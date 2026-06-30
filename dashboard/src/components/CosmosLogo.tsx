import { cn } from "@utils/helpers";
import Image from "next/image";
import * as React from "react";
import CosmosLogoMark from "@/assets/cosmos.svg";
import CosmosLogoFull from "@/assets/cosmos-full.svg";

type Props = {
  size?: "default" | "large";
  mobile?: boolean;
};

const sizes = {
  default: {
    desktop: 30,
    mobile: 30,
  },
  large: {
    desktop: 36,
    mobile: 40,
  },
};

export const CosmosLogo = ({ size = "default", mobile = true }: Props) => {
  return (
    <>
      <Image
        src={CosmosLogoFull}
        height={sizes[size].desktop}
        alt={"Cosmos"}
        className={cn(mobile && "hidden md:block")}
      />
      {mobile && (
        <Image
          src={CosmosLogoMark}
          width={sizes[size].mobile}
          alt={"Cosmos"}
          className={cn("md:hidden ml-4")}
        />
      )}
    </>
  );
};
