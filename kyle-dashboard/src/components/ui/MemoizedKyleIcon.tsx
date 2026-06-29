import * as React from "react";
import { memo } from "react";
import KyleIcon from "@/assets/icons/KyleIcon";

const MemoizedKyleIcon = () => {
  return <KyleIcon size={14} />;
};

export default memo(MemoizedKyleIcon);
