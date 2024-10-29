import { ComponentType } from 'react';

export function withEdition<T>(
  WrappedComponent: ComponentType<T>,
  edition: 'BE'
): ComponentType<T> {
  // Try to create a nice displayName for React Dev Tools.
  const displayName =
    WrappedComponent.displayName || WrappedComponent.name || 'Component';

  function WrapperComponent(props: T & JSX.IntrinsicAttributes) {
    return <WrappedComponent {...props} />;
  }

  WrapperComponent.displayName = `with${edition}Edition(${displayName})`;

  return WrapperComponent;
}
