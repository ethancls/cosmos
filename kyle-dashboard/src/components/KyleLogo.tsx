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
    desktop: 22,
    mobile: 30,
  },
  large: {
    desktop: 24,
    mobile: 40,
  },
};

export const KyleLogo = ({ size = "default", mobile = true }: Props) => {
  return (
    <div className="flex items-center gap-2">
      <Image
        src={KyleLogoMark}
        height={sizes[size].desktop}
        alt={"Kyle Logo"}
        className={cn(mobile && "hidden md:block")}
      />
      {mobile && (
        <Image
          src={KyleLogoMark}
          width={sizes[size].mobile}
          alt={"Kyle Logo"}
          className={cn(mobile && "md:hidden")}
        />
      )}
      <span className="text-lg font-bold text-gray-900 dark:text-white hidden md:block">
        Kyle
      </span>
    </div>
  );
};
