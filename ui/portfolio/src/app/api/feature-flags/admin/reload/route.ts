import { NextRequest, NextResponse } from "next/server";

const FEATURE_FLAGS_API_URL =
  process.env.FEATURE_FLAGS_API_URL || "http://localhost:4000";

export async function POST(request: NextRequest) {
  try {
    const response = await fetch(`${FEATURE_FLAGS_API_URL}/admin/reload`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const errorData = await response.json();
      return NextResponse.json(
        { error: errorData.error || "Failed to reload flags" },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error reloading flags:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    );
  }
}
