import { NextRequest, NextResponse } from "next/server";

const OUTBOX_API_URL =
  process.env.NEXT_PUBLIC_OUTBOX_API_URL || "http://localhost:8080";

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;

    const response = await fetch(
      `${OUTBOX_API_URL}/api/v1/events/${id}/retry`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      if (response.status === 404) {
        return NextResponse.json({ error: "Event not found" }, { status: 404 });
      }
      if (response.status === 400) {
        const errorData = await response.json();
        return NextResponse.json(errorData, { status: 400 });
      }
      throw new Error(`Outbox API error: ${response.status}`);
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error retrying event:", error);
    return NextResponse.json(
      { error: "Failed to retry event" },
      { status: 500 }
    );
  }
}
