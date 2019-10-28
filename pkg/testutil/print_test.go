package testutil

import (
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	mockexpfmt "github.com/nieltg/prom-example-testutil/test/mock_expfmt"
	mocktestutil "github.com/nieltg/prom-example-testutil/test/mock_testutil"
	prommodel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

var printMetricsNameA = "name-a"
var printMetricsNameB = "name-b"
var printMetricsA = []*prommodel.MetricFamily{
	&prommodel.MetricFamily{Name: &filterMetricsNameA},
}
var printMetricsB = []*prommodel.MetricFamily{
	&prommodel.MetricFamily{Name: &filterMetricsNameA},
	&prommodel.MetricFamily{Name: &filterMetricsNameB},
}
var errPrintA = fmt.Errorf("error-a")

func mockPrinter(printer printer) func() {
	originalPrinter := globalPrinter
	globalPrinter = printer

	return func() {
		globalPrinter = originalPrinter
	}
}

func TestMustPrintMetrics(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	printer := mocktestutil.NewMockprinter(controller)
	printer.EXPECT().PrintMetrics(printMetricsA)

	defer mockPrinter(printer)()
	MustPrintMetrics(printMetricsA)
}

func TestMustPrintMetrics_panic(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	printer := mocktestutil.NewMockprinter(controller)
	printer.EXPECT().PrintMetrics(gomock.Any()).Return(errPrintA).AnyTimes()

	defer mockPrinter(printer)()
	assert.PanicsWithValue(t, errPrintA, func() {
		MustPrintMetrics(nil)
	})
}

func TestPrintMetrics(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	printer := mocktestutil.NewMockprinter(controller)
	printer.EXPECT().PrintMetrics(printMetricsA).Return(errPrintA)

	defer mockPrinter(printer)()
	err := PrintMetrics(printMetricsA)
	t.Run("error", func(t *testing.T) {
		assert.EqualError(t, err, errPrintA.Error())
	})
}

func newPrinterWithEncoder(encoder expfmt.Encoder) printer {
	return &printerImpl{
		newEncoderFunc: func(w io.Writer, format expfmt.Format) expfmt.Encoder {
			return encoder
		},
	}
}

func Test_printImpl_PrintMetrics(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockEncoder := mockexpfmt.NewMockEncoder(controller)
	mockEncoder.EXPECT().Encode(printMetricsA[0]).Return(nil).Times(1)

	printer := newPrinterWithEncoder(mockEncoder)
	_ = printer.PrintMetrics(printMetricsA)
}

func Test_printImpl_PrintMetrics_nil(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockEncoder := mockexpfmt.NewMockEncoder(controller)
	mockEncoder.EXPECT().Encode(gomock.Any()).Return(nil).Times(0)

	printer := newPrinterWithEncoder(mockEncoder)
	_ = printer.PrintMetrics(nil)
}

func Test_printImpl_PrintMetrics_multiple(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockEncoder := mockexpfmt.NewMockEncoder(controller)
	gomock.InOrder(
		mockEncoder.EXPECT().Encode(printMetricsB[0]),
		mockEncoder.EXPECT().Encode(printMetricsB[1]),
	)

	printer := newPrinterWithEncoder(mockEncoder)
	_ = printer.PrintMetrics(printMetricsB)
}

func Test_printImpl_PrintMetrics_error(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockEncoder := mockexpfmt.NewMockEncoder(controller)
	mockEncoder.EXPECT().Encode(gomock.Any()).Return(errPrintA).AnyTimes()

	printer := newPrinterWithEncoder(mockEncoder)
	assert.EqualError(t, printer.PrintMetrics(printMetricsA), errPrintA.Error())
}
