import { NextRequest, NextResponse } from "next/server";

const FEATURE_FLAGS_API_URL =
  process.env.FEATURE_FLAGS_API_URL || "http://localhost:4000";

export async function GET(request: NextRequest) {
  try {
    const searchParams = request.nextUrl.searchParams;
    const env = searchParams.get("env") || "local";

    const response = await fetch(`${FEATURE_FLAGS_API_URL}/flags?env=${env}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const errorData = await response.json();
      return NextResponse.json(
        { error: errorData.error || "Failed to fetch flags" },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error fetching flags:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    );
  }
}
