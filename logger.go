package ctxLogger

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var TraceLog *zap.Logger

const traceId = "trace_id"

func init() {
	content, _ := rotatelogs.New(
		"./logs/%Y%m%d/info.log"+"-%Y%m%d%H%M",
		rotatelogs.WithRotationTime(4*time.Hour), //rotate 最小为1分钟轮询。默认60s  低于1分钟就按1分钟来
		rotatelogs.WithMaxAge(time.Hour*24*30),   //日志保存默认保存一个月吧
	)
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	consoleEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	deBugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	consoleDebugging := zapcore.Lock(os.Stdout)
	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := content
	//errorWriter := content

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleDebugging, deBugLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(infoWriter), infoLevel), //级别比info大的都输出到文件中
	)

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	TraceLog = log
}
func Info(ctx *gin.Context, template string, fields ...zap.Field) {
	if ctx != nil {
		TraceLog.With(fields...).Info(template, zap.String("trace_id", ctx.Request.Header.Get(traceId)))
	} else {
		TraceLog.With(fields...).Info(template)
	}
}
func Warn(ctx *gin.Context, template string, fields ...zap.Field) {
	if ctx != nil {
		TraceLog.With(fields...).Warn(template, zap.String("trace_id", ctx.Request.Header.Get(traceId)))
	} else {
		TraceLog.With(fields...).Warn(template)
	}

}
func Debug(ctx *gin.Context, template string, fields ...zap.Field) {
	if ctx != nil {
		TraceLog.With(fields...).Debug(template, zap.String("trace_id", ctx.Request.Header.Get(traceId)))
	} else {
		TraceLog.With(fields...).Debug(template)
	}
}
func Error(ctx *gin.Context, template string, fields ...zap.Field) {
	if ctx != nil {
		TraceLog.With(fields...).Error(template, zap.String("trace_id", ctx.Request.Header.Get(traceId)))
	} else {
		TraceLog.With(fields...).Error(template)
	}
}

//func for fiber
func FInfo(ctx *fiber.Ctx, template string, fields ...zap.Field) {
	if ctx != nil {
		TraceLog.With(fields...).Info(template, zap.ByteString("trace_id", ctx.Response().Header.Peek(fiber.HeaderXRequestID)))
	} else {
		TraceLog.With(fields...).Info(template)
	}
}
func FWarn(ctx *fiber.Ctx, template string, fields ...zap.Field) {
	if ctx != nil {
		TraceLog.With(fields...).Warn(template, zap.ByteString("trace_id", ctx.Response().Header.Peek(fiber.HeaderXRequestID)))
	} else {
		TraceLog.With(fields...).Warn(template)
	}

}
func FDebug(ctx *fiber.Ctx, template string, fields ...zap.Field) {
	if ctx != nil {
		TraceLog.With(fields...).Debug(template, zap.ByteString("trace_id", ctx.Response().Header.Peek(fiber.HeaderXRequestID)))
	} else {
		TraceLog.With(fields...).Debug(template)
	}
}
func FError(ctx *fiber.Ctx, template string, fields ...zap.Field) {
	if ctx != nil {
		TraceLog.With(fields...).Error(template, zap.ByteString("trace_id", ctx.Response().Header.Peek(fiber.HeaderXRequestID)))
	} else {
		TraceLog.With(fields...).Error(template)
	}
}
