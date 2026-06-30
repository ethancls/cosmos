import { cn, generateColorFromString, getGravatarUrl } from "@utils/helpers";
import { Cog } from "lucide-react";
import * as React from "react";
import { useEffect, useState } from "react";
import Image from "next/image";

type Props = {
  name?: string;
  email?: string;
  id?: string;
  picture?: string;
  size?: "default" | "sm";
  className?: string;
};
export const SmallUserAvatar = ({
  name,
  id,
  email,
  picture,
  size = "default",
  className,
}: Props) => {
  const [pictureFailed, setPictureFailed] = useState(false);
  const [gravatarFailed, setGravatarFailed] = useState(false);
  const [gravatarUrl, setGravatarUrl] = useState<string | undefined>();

  const isCosmos = email === "Cosmos";
  const px = size === "sm" ? 20 : 28;

  useEffect(() => {
    getGravatarUrl(email, px).then(setGravatarUrl);
  }, [email, px]);

  // 1. Profile picture
  if (!pictureFailed && picture) {
    return (
      <Image
        src={picture}
        alt={name ?? ""}
        onError={() => setPictureFailed(true)}
        width={px}
        height={px}
        className={cn("rounded-full shrink-0", className)}
      />
    );
  }

  // 2. Gravatar
  if (!gravatarFailed && gravatarUrl) {
    return (
      <Image
        src={gravatarUrl}
        alt={name ?? ""}
        onError={() => setGravatarFailed(true)}
        width={px}
        height={px}
        className={cn("rounded-full shrink-0", className)}
      />
    );
  }

  // 3. Letter avatar fallback
  return (
    <div
      className={cn(
        "rounded-full shrink-0 flex items-center justify-center text-white uppercase font-medium bg-nb-gray-850",
        size === "default" && "w-7 h-7 text-[12px]",
        size === "sm" && "w-5 h-5 text-[9px] leading-[0]",
        className,
      )}
      style={{
        color: isCosmos
          ? "#808080"
          : generateColorFromString(name || id || "System User"),
      }}
    >
      {isCosmos ? (
        <Cog size={14} />
      ) : (
        name?.charAt(0) || id?.charAt(0)
      )}
    </div>
  );
};
