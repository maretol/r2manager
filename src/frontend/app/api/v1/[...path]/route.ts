import { NextRequest, NextResponse } from 'next/server'

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8080'

async function proxyRequest(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }): Promise<NextResponse> {
  const { path } = await params
  const pathString = path.map(encodeURIComponent).join('/')
  const searchParams = request.nextUrl.searchParams.toString()
  const url = `${BACKEND_URL}/api/v1/${pathString}${searchParams ? `?${searchParams}` : ''}`

  const headers = new Headers(request.headers)
  headers.delete('host')

  const init: RequestInit = {
    method: request.method,
    headers,
  }

  if (request.body) {
    init.body = request.body
    // @ts-expect-error duplex is needed for streaming request bodies
    init.duplex = 'half'
  }

  const response = await fetch(url, init)

  const responseHeaders = new Headers(response.headers)
  responseHeaders.delete('transfer-encoding')

  return new NextResponse(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: responseHeaders,
  })
}

export const GET = proxyRequest
export const POST = proxyRequest
export const PUT = proxyRequest
export const DELETE = proxyRequest
export const PATCH = proxyRequest
