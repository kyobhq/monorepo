const VIDEO_EXTENSIONS = ['mp4'];

export function getFileType(url: string): 'image' | 'video' | 'unknown' {
  const split = url.split('/');
  const fileName = split[split.length - 1];
  const [_, extension] = fileName.split('.');

  if (extension === 'webp') return 'image';
  if (VIDEO_EXTENSIONS.includes(extension)) return 'video';

  return 'unknown';
}

export function getOriginalFileName(url: string): string {
  const split = url.split('/');
  const fileName = split[split.length - 1];
  const splitFileName = fileName.split('-');

  return splitFileName[1];
}
