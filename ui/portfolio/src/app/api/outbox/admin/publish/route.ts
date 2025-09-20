import { NextRequest, NextResponse } from "next/server";

const OUTBOX_API_URL =
  process.env.NEXT_PUBLIC_OUTBOX_API_URL || "http://localhost:8080";

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    const response = await fetch(`${OUTBOX_API_URL}/admin/publish`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });

    if (!response.ok) {
      throw new Error(`Outbox API error: ${response.status}`);
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error publishing events:", error);
    return NextResponse.json(
      { error: "Failed to publish events" },
      { status: 500 }
    );
  }
}
