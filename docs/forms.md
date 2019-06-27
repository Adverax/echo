# Концепция
Формы предназначены для отображения и правки клиентом каких-либо данных. Однако, чаще всего одни и те же данные, хранимые в базе данных и отображаемые клиенту, имеют различные представления. Например, предоставляя возможность выбора языка клиенту, мы обычно выводим названия этих языков, в то время, как в базе данных хранится только идентификатор языка. Таким образом, у нас возникает необходимость хранить данные во внутреннем представлении, а отображать - во внешнем представлении. Как правило, внешнее представление имеет всегда литеральный формат. Для перекодировки данных используются кодеки. 

# Общая схема
Общая схема написания обработчика обновления выглядит следующим образом:
* запрашиваем редактируемый объект
* определяем поля
* на основании полей формируем модель
* принимаем решение о корректности данных
* если модель корректна, то обрабатываем полученные данные, после чего перенаправляем клиента на новую страницу.
* если модель недостоверна, то формируем страницу.

```go
package user

import (
	"github.com/adverax/echo"
	"github.com/adverax/echo/widget"
	"net/http"
)

func actionUpdate(ctx echo.Context) error {
	// Get current value of editable object
	obj := getMyObj(ctx)
	
	// Declare fields
	field1 := ...
	field2 := ...
	
	// Make model from declared fields
    m := echo.Model{
        "Field1": field1,
        "Field2": field2,
    }

    // Resolve model (import, validate and export)
    if err := m.Resolve(ctx, obj, obj, nil); err != nil {
        if err != echo.ErrModelSealed {
        	// Unknown error 
            return err
        }

        // Model is valid and data exported into the editable object
        // ...

        // Redirect to the target url 
        return ctx.Redirect(http.StatusSeeOther, "/user/view")
    }

    m["Submit"] = &widget.FormSubmit{
        Label: "Update",
    }

    // Model is invalid
    content := widget.Map{
        "Form": &widget.Form{
            Name: "update-form",
            Model: m,
        },
    }
    
    // Make response from content
    // ...
    return nil
}
```

# Импорт экспорт данных модели
Каждая модель может импортировать и экспортировать данные из/в изменяемую структуру данных. Для этой цели существует пара методов
* Import(ctx echo.Context, src interface{}, mapper echo.Mapper) error 
* Export(ctx echo.Context, dst interface{}, mapper echo.Mapper) error
Параметр mapper служит для преобразования имен между структурой данных и моделью. Здесь ключ - название полея в модели, а значение - название поля в структуре. Причем, название поля может быть составным, например: "Coordinate.Latitude". Это позволяет использовать вложенные структуры данных. Все имена используемых полей структур данных должны начинаться с заглавных букв, поскольку они должны быть экспортируемыми.

# Компоненты
Определения компонентов позволяют использовать одну и ту же информацию, как на стороне сервера, так и на стороне клиента.

## Общие свойства компонентов
* Label - подпись виджета, которая отображается клиенту для пояснения назначения этого компонента.
* Default - значение по умолчанию, которое будет принимать компонент без какого-либо участия со стороны клиента.
* Required - обязательное поле, которое требует обязательного заполнения.
* Pattern - шаблон регулярного выражения, которому должно соответствовать внешнее представление.
* Filter - функция фильтрации (предварительной обработки) литерального представления. Например, используя функцию strings.TrimSpaces мы можем удалять начальные и хвостовые пробелы.
* Disabled - поля, доступ к которым запрещен (в том числе и их изменение).

## Перечень компонентов
* FormText - текстовое поле, которое может быть использовано как для однострочного, так и для многострочного редактора.
* FormSelect - используется для создания html combobox.
* FormFlag - используется для создания независимого ckeckbox.
* FormFlags - используется для создания зависимого radio.
* FormSubmit - испольуется для создания как одиночного submit, так и для группы элементов submit.
* FormHidden - используется для создания скрытого поля hidden.
* FormFile - используется для создания html file.

# Валидация данных
Любые присылаемые клиентом данные не вызывают доверия и должны быть проверены. Для их проверки используется механизм валидации. При возникновении ошибки валидации, операция не может быть завершена и клиенту предлагается исправить возникшую ошибку. Большинство ошибок валидации генерируется непосредственно самими компонентами или же кодеками. При необходимости, можно выполнить также позднюю проверку данных и сгенерировать требуемые ошибки. Например:

```go
package user

import (
	"github.com/adverax/echo"
	"github.com/adverax/echo/widget"
	"net/http"
)

type Range struct {
	Min uint8
	Max uint8
}

func actionUpdate(ctx echo.Context) error {
	min := &widget.FormText{
		Name: "min",
		Codec: echo.Uint8Codec,
	}

	max := &widget.FormText{
		Name: "Max",
		Codec: echo.Uint8Codec,
	}

	// Make model
	m := &echo.Model{
		"Min": min,
		"Max": max,
	} 
	
	var r Range
    	
    // Resolve model (import, validate and export)
    if err := m.Resolve(ctx, nil, &r, nil); err != nil {
        if err != echo.ErrModelSealed {
        	// Unknown error 
            return err
        }
        
        // Do late validation
        if r.Min > r.Max {
        	min.AddError(echo.NewValidationErrorString("Value too large"))
        	max.AddError(echo.NewValidationErrorString("Value too small"))
        	return nil
        }

        // Model is valid and data exported into the editable object
        // ...

        // Redirect to the target url 
        return ctx.Redirect(http.StatusSeeOther, "/user/view")
    }

    m["Submit"] = &widget.FormSubmit{
        Label: "Update",
    }

    // Model is invalid
    content := widget.Map{
        "Form": &widget.Form{
            Name: "update-form",
            Model: m,
        },
    }
    
    // Make response from content
    // ...
    return nil
}
``` 
 
# Множественные формы
Иногда бывает необходимо использовать несколько форм на одной странице для единовремненой правки данных. Типичным примером такого подхода является форма правки статьи сразу на всех доступных языках. В результате мы получаем страницу, содержащую ленту форм.

Процесс взаимодействия с такой множественной формой почти не отличается от обычных форм, за исключением того, что свойство Models содержит последовательность необходимых форм. Такая комплексная форма при валидации выполняет проверку всех входящих в нее форм.

# Многошаговые формы
Многошаговые формы - это формы, расположенные на ряде веб страниц и хранящая свое состояние (MultiStepState) в кэше.

Координатором верхнего уровня такой формы выступает стратегия:

```go
type MultiStepStrategy interface {
	// Started initialization. Executed for restart strategy only.
	Setup(ctx echo.Context, values generic.Params) (stage string, err error)
	// Create stage instance
	Stage(ctx echo.Context, state *MultiStepState) (MultiStepStage, error)
}
```

Здесь метод Setup обеспечивает первоначальную настройку параметров и возвращает название начального шага. Каждый шаг должен реализовывать следующий интерфейс:

```go
type MultiStepStage interface {
	// Get model of stage
	Model(ctx echo.Context, state *MultiStepState) (echo.Model, error)
	// Import data into the model and validate it.
	Resolve(ctx echo.Context, state *MultiStepState, model echo.Model) error
	// Consume stage data. Return new state reference or nil (see method MultiStepState.Become)
	Consume(ctx echo.Context, state *MultiStepState, model echo.Model) (reply interface{}, err error)
	// Publish form
	Publish(ctx echo.Context, state *MultiStepState, form *MultiStepForm) error
}
``` 

Для получения экземпляра шага по его имени служит метод MultiStepStrategy.Stage.

Для упрощения реализации, чаще всего шаг реализуется на основании прототипа MultiStepBaseStage, реализующего базовые методы Resolve и Publish:

```go
type MyStage struct {
	MultiStepBaseStage
}

func (stage *MyStage) Model(ctx echo.Context, state *MultiStepState) (echo.Model, error) {
	
}

func (stage *MyStage) Consume(ctx echo.Context, state *MultiStepState, model echo.Model) (reply interface{}, err error) {
	
}
```

Как понятно из названия метода, Model служит построения модели шага. Метод же Consume, служит для потребления поступивших данных модели (по необходимости) и должен возвращать адрес нового шага:
* имя шага - переход на этот шаг.
* число - служит для возврата на указанное число шагов назад по истории навигации.
* nil - тупик и служит для завершения работы многошаговой формы.

Для подключения такой формы к роутеру, должна использоваться функция InitMultiStepRouter, обеспечивающая подключение нескольких обработчиков к роутеру.
