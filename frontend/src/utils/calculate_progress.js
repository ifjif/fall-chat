export const calculateProgress = (part, total) => {
  const progress = Math.ceil(part * (total / 100))
  if (progress >= 100) return 99
  return progress
}
