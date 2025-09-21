import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  try {
    const event = await request.json();

    return NextResponse.json({
      message: "Event received successfully",
      eventId: event.id,
    });
  } catch (error) {
    console.error("Error processing webhook:", error);
    return NextResponse.json(
      { error: "Failed to process webhook" },
      { status: 500 }
    );
  }
}

export async function GET() {
  return NextResponse.json({
    message: "Webhook endpoint is active",
    timestamp: new Date().toISOString(),
  });
}
