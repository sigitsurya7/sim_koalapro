import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const PUBLIC_PATHS = ['/login', '/auth/login', '/koala-favicon.ico']

const isPublicPath = (pathname: string) => {
  if (PUBLIC_PATHS.includes(pathname)) return true
  return pathname.startsWith('/_next') || pathname.startsWith('/api/public') || pathname.startsWith('/static')
}

export function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl
  const token = req.cookies.get('token')?.value

  if (!token && !isPublicPath(pathname)) {
    const loginUrl = req.nextUrl.clone()
    loginUrl.pathname = '/login'
    loginUrl.searchParams.set('from', pathname)
    return NextResponse.redirect(loginUrl)
  }

  if (token && pathname === '/login') {
    const dashboardUrl = req.nextUrl.clone()
    dashboardUrl.pathname = '/dashboard'
    return NextResponse.redirect(dashboardUrl)
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!api).*)'],
}
