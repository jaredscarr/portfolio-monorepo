import { NextResponse } from "next/server";

const OBSERVABILITY_API_URL =
  process.env.OBSERVABILITY_API_URL || "http://localhost:8081";

export async function GET() {
  try {
    const response = await fetch(`${OBSERVABILITY_API_URL}/health`, {
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
    console.error("Error checking observability health:", error);
    return NextResponse.json({ status: "unhealthy" }, { status: 503 });
  }
}
