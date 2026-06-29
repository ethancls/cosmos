import { cn } from "@utils/helpers";
import Image from "next/image";
import * as React from "react";
import KyleLogoMark from "@/assets/kyle.svg";
import KyleLogoFull from "@/assets/kyle-full.svg";

type Props = {
  size?: "default" | "large";
  mobile?: boolean;
};

const sizes = {
  default: {
    desktop: 24,
    mobile: 30,
  },
  large: {
    desktop: 32,
    mobile: 40,
  },
};

export const KyleLogo = ({ size = "default", mobile = true }: Props) => {
  return (
    <>
      <Image
        src={KyleLogoFull}
        height={sizes[size].desktop}
        alt={"Kyle"}
        className={cn(mobile && "hidden md:block")}
      />
      {mobile && (
        <Image
          src={KyleLogoMark}
          width={sizes[size].mobile}
          alt={"Kyle"}
          className={cn("md:hidden ml-4")}
        />
      )}
    </>
  );
};
