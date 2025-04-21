// next.config.ts for Next.js (if you're using TypeScript)

import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  env: {
    NEXT_PUBLIC_API_PATH: process.env.NEXT_PUBLIC_API_PATH,  // Access API_PATH from .env
  },
};

export default nextConfig;
