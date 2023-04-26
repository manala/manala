export function ltrim(str: string, charlist: string) {
  return str.replace(new RegExp(`^[${charlist}]+`), '');
}

export function rtrim(str: string, charlist: string) {
  return str.replace(new RegExp(`[${charlist}]+$`), '');
}
