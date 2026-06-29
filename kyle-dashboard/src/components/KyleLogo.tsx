import { cn } from "@utils/helpers";
import Image from "next/image";
import * as React from "react";
import KyleLogoMark from "@/assets/kyle.svg";

type Props = {
  size?: "default" | "large";
  mobile?: boolean;
};

const sizes = {
  default: {
    desktop: 28,
    mobile: 36,
  },
  large: {
    desktop: 32,
    mobile: 44,
  },
};

export const KyleLogo = ({ size = "default", mobile = true }: Props) => {
  return (
    <div className="flex items-center gap-2.5">
      <Image
        src={KyleLogoMark}
        height={sizes[size].desktop}
        width={sizes[size].desktop}
        alt={"Kyle Logo"}
        className={cn(mobile && "hidden md:block")}
      />
      {mobile && (
        <Image
          src={KyleLogoMark}
          width={sizes[size].mobile}
          height={sizes[size].mobile}
          alt={"Kyle Logo"}
          className={cn("md:hidden")}
        />
      )}
      <span className="text-xl font-medium text-white hidden md:block">
        Kyle
      </span>
    </div>
  );
};
