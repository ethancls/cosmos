import { cn, generateColorFromUser, getGravatarUrl } from "@utils/helpers";
import * as React from "react";
import { useEffect, useState } from "react";
import Image from "next/image";
import { useApplicationContext } from "@/contexts/ApplicationProvider";

type Props = {
  size?: "default" | "small" | "large" | "medium";
};

export const UserAvatar = ({ size = "default" }: Props) => {
  const { user } = useApplicationContext();

  const [pictureFailed, setPictureFailed] = useState(false);
  const [gravatarFailed, setGravatarFailed] = useState(false);
  const [gravatarUrl, setGravatarUrl] = useState<string | undefined>();

  const getAvatarSize = () => {
    if (size === "small") return 32;
    if (size === "default") return 40;
    if (size === "large") return 48;
    return 35.2;
  };

  const sizePx = getAvatarSize();

  useEffect(() => {
    getGravatarUrl(user?.email, sizePx).then(setGravatarUrl);
  }, [user?.email, sizePx]);

  // 1. Profile picture from backend
  if (!pictureFailed && user?.picture) {
    return (
      <Image
        src={user.picture}
        alt={user.name ?? ""}
        onError={() => setPictureFailed(true)}
        width={sizePx}
        height={sizePx}
        className={"rounded-full"}
      />
    );
  }

  // 2. Gravatar
  if (!gravatarFailed && gravatarUrl) {
    return (
      <Image
        src={gravatarUrl}
        alt={user?.name ?? ""}
        onError={() => setGravatarFailed(true)}
        width={sizePx}
        height={sizePx}
        className={"rounded-full"}
      />
    );
  }

  // 3. Letter avatar fallback
  return (
    <div
      className={cn(
        "rounded-full flex items-center justify-center bg-nb-gray-900 uppercase",
        size == "small" && "w-8 h-8",
        size == "medium" && "w-[2.2rem] h-[2.2rem]",
        size == "default" && "w-10 h-10",
        size == "large" && "w-12 h-12",
      )}
      style={{
        color: generateColorFromUser(user),
      }}
    >
      {user?.name?.charAt(0) || user?.id?.charAt(0)}
    </div>
  );
};
