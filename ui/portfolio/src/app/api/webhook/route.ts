import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  try {
    const event = await request.json();

    // Log the received event (in production, you'd want to process this properly)
    console.log("Received webhook event:", {
      id: event.id,
      type: event.type,
      source: event.source,
      timestamp: new Date().toISOString(),
    });

    // Here you would typically:
    // 1. Validate the event
    // 2. Process the event data
    // 3. Store in your database
    // 4. Trigger any necessary business logic

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
