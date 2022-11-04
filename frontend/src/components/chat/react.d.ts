interface RefObject<T> {
  // immutable
  current: T | null;
  style: T | null;
}

function createRef<T>(): RefObject<T>
