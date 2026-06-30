import * as React from "react";
import { memo } from "react";
import CosmosIcon from "@/assets/icons/CosmosIcon";

const MemoizedCosmosIcon = () => {
  return <CosmosIcon size={14} />;
};

export default memo(MemoizedCosmosIcon);
