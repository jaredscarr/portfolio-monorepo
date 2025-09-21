import { NextRequest, NextResponse } from "next/server";

const FEATURE_FLAGS_API_URL = "http://localhost:4000";

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ key: string }> }
) {
  try {
    const { searchParams } = new URL(request.url);
    const env = searchParams.get("env") || "local";
    const body = await request.json();
    const { key } = await params;

    const response = await fetch(
      `${FEATURE_FLAGS_API_URL}/admin/flags/${key}?env=${env}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      }
    );

    if (!response.ok) {
      throw new Error(`Feature flags API error: ${response.status}`);
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error updating feature flag:", error);
    return NextResponse.json(
      { error: "Failed to update feature flag" },
      { status: 500 }
    );
  }
}
