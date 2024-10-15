import { toast } from '@/components/ui/toast'

export function errorToast(message: string) {
    toast({
        title: 'An error happened',
        variant: 'destructive',
        description: message
    })
}
