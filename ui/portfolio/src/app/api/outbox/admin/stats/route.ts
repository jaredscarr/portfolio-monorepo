import { NextResponse } from "next/server";

const OUTBOX_API_URL =
  process.env.NEXT_PUBLIC_OUTBOX_API_URL || "http://localhost:8080";

export async function GET() {
  try {
    const response = await fetch(`${OUTBOX_API_URL}/admin/stats`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Outbox API error: ${response.status}`);
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error fetching stats:", error);
    return NextResponse.json(
      { error: "Failed to fetch stats" },
      { status: 500 }
    );
  }
}
