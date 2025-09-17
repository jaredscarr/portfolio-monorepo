import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  turbopack: {
    resolveAlias: {
      // Add any specific aliases if needed
    },
    resolveExtensions: [".js", ".jsx", ".ts", ".tsx"],
  },
};

export default nextConfig;
