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
    mobile: 34,
  },
  large: {
    desktop: 36,
    mobile: 44,
  },
};

export const KyleLogo = ({ size = "default", mobile = true }: Props) => {
  return (
    <div className="flex items-center gap-3">
      <Image
        src={KyleLogoMark}
        width={sizes[size].desktop}
        height={sizes[size].desktop}
        alt={"Kyle"}
        className={cn("shrink-0", mobile && "hidden md:block")}
      />
      {mobile && (
        <Image
          src={KyleLogoMark}
          width={sizes[size].mobile}
          height={sizes[size].mobile}
          alt={"Kyle"}
          className={cn("shrink-0 md:hidden")}
        />
      )}
      <span className="text-[22px] font-semibold tracking-wide text-white hidden md:block">
        Kyle
      </span>
    </div>
  );
};
