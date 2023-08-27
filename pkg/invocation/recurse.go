package invocation

// type Invokable[I, O any] interface {
// 	// Handler[I, O]
// 	Invoke(ctx context.Context, input I) (O, error)
// }

// type LambdaInvoker[I, O any] struct {
// 	invoker Invoker[I, O]
// }

// // func (me LambdaInvoker[I, O]) Invoke(ctx context.Context, raw []byte) ([]byte, error) {

// // 	me.invoker.Logger().Debug().Str("raw", string(raw)).Msg("Invoking Lambda")

// // 	var in I
// // 	err := xon.Unmarshal(raw, &in)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	out, err := me.invoker.Invoke(ctx, in)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	return xon.Marshal(out)
// // }

// func NewLambdaInvoker[I any, O any](me Invoker[I, O]) LambdaInvoker[I, O] {
// 	return LambdaInvoker[I, O]{me}
// }

// func (me LambdaInvoker[I, O]) WrapWithTracing() func(context.Context, interface{}) (interface{}, error) {

// 	tp, err := xrayconfig.NewTracerProvider(me.invoker.RootContext())
// 	if err != nil {
// 		me.invoker.Logger().Error().Err(err).Msg("error creating tracer provider")
// 	}

// 	me.invoker.AddCleanup(func(Handler[I, O]) {
// 		err := tp.Shutdown(me.invoker.RootContext())
// 		if err != nil {
// 			me.invoker.Logger().Error().Err(err).Msg("error shutting down tracer provider")
// 		}
// 	})

// 	otel.SetTracerProvider(tp)
// 	otel.SetTextMapPropagator(xray.Propagator{})

// 	wrapped := otellambda.InstrumentHandler(me.Invoke, xrayconfig.WithRecommendedOptions(tp)...)

// 	if wrp, ok := wrapped.(func(context.Context, interface{}) (interface{}, error)); ok {
// 		return wrp
// 	} else {
// 		me.invoker.Logger().Error().Interface("wrapped", wrapped).Msg("wrapped handler is not a function")
// 		panic("wrapped handler is not a function")
// 	}

// }
