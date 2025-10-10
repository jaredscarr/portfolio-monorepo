import { NextResponse } from "next/server";

const FEATURE_FLAGS_API_URL =
  process.env.FEATURE_FLAGS_API_URL || "http://localhost:4000";

export async function GET() {
  try {
    const response = await fetch(`${FEATURE_FLAGS_API_URL}/health`, {
      method: "GET",
    });

    if (!response.ok) {
      return NextResponse.json(
        { status: "unhealthy" },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error checking feature-flags health:", error);
    return NextResponse.json({ status: "unhealthy" }, { status: 503 });
  }
}
