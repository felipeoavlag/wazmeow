"use client";

import { useEffect } from "react";
import { useUIStore } from "@/lib/stores/app-store";

export function useMobile() {
  const { isMobile, setIsMobile } = useUIStore();

  useEffect(() => {
    const checkMobile = () => {
      const mobile = window.innerWidth < 768; // md breakpoint
      setIsMobile(mobile);
    };

    // Check on mount
    checkMobile();

    // Add event listener
    window.addEventListener("resize", checkMobile);

    // Cleanup
    return () => window.removeEventListener("resize", checkMobile);
  }, [setIsMobile]);

  return isMobile;
}