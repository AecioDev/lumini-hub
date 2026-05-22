import { Building2 } from "lucide-react"

interface LogoProps {
  className?: string
}

export function Logo({ className }: LogoProps) {
  return (
    <div className={`flex items-center justify-center rounded-md bg-primary p-2 text-primary-foreground ${className}`}>
      <Building2 className="h-8 w-8" />
    </div>
  )
}

