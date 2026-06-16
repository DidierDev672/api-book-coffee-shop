# Paulu — Guía de uso del panel de control

> *"Un guerrero no espera órdenes. Lidera desde el frente."* — Paulus Atreides

**Paulu** es el rostro visible de tu sistema Atreides. Así como el Duque Paulus Atreides lideraba a sus soldados desde el frente — carismático, valiente, impulsivo pero siempre fiel a los valores de su Casa — esta interfaz es tu herramienta de batalla diaria: ágil, directa y diseñada para que tomes decisiones en el momento.

Mientras Atreides (el backend) gobierna los datos en las sombras como un Mentat calculando rutas de especias, Paulu es el guerrero que tú ves y con el que interactúas: el panel desde donde registras ventas, apruebas órdenes, controlas tu inventario y lideras tu negocio.

Esta guía te enseñará la operación más importante del día a día: **registrar una salida de inventario (venta) y gestionar tus productos.**

---

## Requisitos previos

Antes de comenzar, asegúrate de tener lo siguiente:

- Una **cuenta activa** en Atreides (ver guía de registro si eres nuevo)
- Tu **empresa registrada** en el sistema
- Al menos **un producto creado** en tu inventario
- Una **bodega o almacén** registrado
- (Opcional) Un **proveedor** registrado si la venta involucra devolución

> **¿Te falta alguno de estos?** Completarlos solo te tomará unos minutos desde el menú lateral. Como decía Paulus: *"Un ejército bien preparado ya tiene media batalla ganada."*

---

## Pasos para registrar tu primera salida (venta)

### Paso 1: Accede al panel de Paulu

*Paulus siempre entraba primero al campo de batalla. Tú entra primero a tu panel.*

1. Abre tu navegador y ve a la dirección de tu sistema.
2. Inicia sesión con tu **correo electrónico** y **contraseña**.
3. Serás recibido por el **Dashboard**, tu centro de mando.

   *[Inserta aquí la captura de pantalla del Dashboard con el menú lateral visible]*

En el menú lateral izquierdo verás varias opciones. Las más importantes para tu día a día son:

| Opción | Para qué sirve |
|--------|----------------|
| **Registrar productos** | Añadir nuevos productos a tu catálogo |
| **Registrar entradas** | Recibir mercancía de proveedores |
| **Registrar orden** | Crear órdenes de compra o reabastecimiento |
| **Registrar salidas** | ** Registrar ventas, devoluciones o ajustes de inventario** |
| **Órdenes** | Ver y gestionar todas tus órdenes |
| **Entradas** | Historial de entradas de producto |

---

### Paso 2: Abre el formulario de salidas

*Un líder carismático sabe exactamente a dónde dirigir a su gente. Tú sabes lo que necesitas hacer.*

1. En el menú lateral, haz clic en **"Registrar salidas"** (icono de exportar).

   *[Inserta aquí la captura de pantalla del menú lateral con "Registrar salidas" resaltado]*

2. Se abrirá el formulario **"Registrar nuevo despacho"**.

   *[Inserta aquí la captura de pantalla del formulario de despacho vacío]*

---

### Paso 3: Completa la información general del despacho

*Paulus no perdía tiempo en rodeos. Iba al grano. Tú también.*

1. **Número de salida** — El sistema genera uno automáticamente con el formato `SAL-YYYYMMDD-XXXX`. Si prefieres uno personalizado, haz clic en **"Generar"** para obtener uno nuevo.

2. **Fecha** — Selecciona la fecha de la salida. Por defecto aparece la fecha de hoy.

3. **Tipo de movimiento** — Selecciona el tipo de salida:

   | Tipo | ¿Cuándo usarlo? |
   |------|------------------|
   | **Venta** | Cuando vendes un producto a un cliente |
   | **Devolución a proveedor** | Cuando devuelves mercancía a tu proveedor |
   | **Donación** | Cuando regalas productos |
   | **Merma** | Cuando productos se dañan o vencen |
   | **Ajuste** | Cuando corriges diferencias de inventario |
   | **Transferencia** | Cuando mueves productos entre bodegas |

   Para una venta normal, selecciona **"Venta"**.

4. **ID de bodega** — Ingresa el identificador de la bodega desde donde sale el producto.

5. **Estado** — Selecciona el estado inicial:

   | Estado | Significado |
   |--------|-------------|
   | **Borrador** | Aún no está confirmado. Puedes editarlo después |
   | **Confirmado** | La salida ya es oficial. El inventario se descuenta |
   | **Cancelado** | La salida no procede |

   Si es una venta que ya realizaste, selecciona **"Confirmado"**. Si es un comprobante que prepararás después, selecciona **"Borrador"**.

   *[Inserta aquí la captura de pantalla del formulario con los campos generales completados]*

---

### Paso 4: Identifica al destinatario

*Todo guerrero sabe para quién lucha. Tú sabes para quién es esta venta.*

1. En la sección **"Destinatario"**, selecciona el **tipo de destinatario**:

   | Tipo | ¿Quién recibe? |
   |------|----------------|
   | **Cliente** | Una persona que compra en tu negocio |
   | **Proveedor** | Un proveedor al que le devuelves mercancía |
   | **Bodega** | Otra bodega tuya (para transferencias) |
   | **Interno** | Consumo interno del negocio |

   Para una venta, selecciona **"Cliente"**. Para una devolución, selecciona **"Proveedor"**.

2. Haz clic en **"Seleccionar"** para elegir el destinatario de tu lista.

   *[Inserta aquí la captura de pantalla del modal de selección de destinatario]*

   > **¿No aparece el destinatario?** Puedes registrarlo desde la opción **"Registrar proveedores"** en el menú lateral.

---

### Paso 5: Agrega los productos a despachar

*El arsenal de Paulus siempre estaba listo. El tuyo también.*

1. En la sección **"Productos a despachar"**, haz clic en **"Agregar producto"**.

   *[Inserta aquí la captura de pantalla del modal de selección de productos]*

2. Se abrirá una ventana con tu lista de productos. Busca por nombre, código o categoría.
3. Selecciona el producto que deseas despachar. Se agregará automáticamente a la tabla.
4. Para cada producto, ajusta:
   - **Cantidad** — Cuántas unidades estás despachando
   - **Costo unitario** — El precio de venta por unidad

   *[Inserta aquí la captura de pantalla de la tabla de productos con cantidades editadas]*

5. Repite los pasos 1–4 para cada producto que incluyas en esta salida.

   > **Consejo:** Paulus era temerario pero no imprudente. Revisa bien las cantidades antes de continuar. Una salida incorrecta afecta tu inventario.

---

### Paso 6: Revisa el resumen financiero

*Hasta el guerrero más valiente debe rendir cuentas a su Casa.*

El sistema calcula automáticamente:

- **Subtotal** — Suma de (cantidad × costo unitario) de cada producto
- **IVA (19%)** — Impuesto al valor agregado
- **Descuento** — Puedes ingresar un descuento si aplica
- **Total** — Subtotal + IVA − Descuento

Puedes agregar **observaciones** en el campo de texto si necesitas dejar notas sobre esta salida (ej: "Venta a crédito", "Factura #1234").

   *[Inserta aquí la captura de pantalla del resumen financiero y observaciones]*

---

### Paso 7: Guarda el despacho

*El momento de la verdad. Paulus no dudaba al dar una orden. Tú tampoco.*

1. Revisa que todos los datos sean correctos.
2. Haz clic en el botón **"Registrar despacho"**.

   *[Inserta aquí la captura de pantalla del botón de guardar]*

3. Si todo está bien, verás un mensaje de éxito:

   | Situación | Mensaje sugerido |
   |-----------|------------------|
   | Despacho registrado | ✅ **"¡Despacho registrado!"** |
   | | *"Todo en orden. El despacho ha quedado registrado en el sistema."* |

4. Decide qué hacer a continuación:
   - **"Crear otro despacho"** — Para registrar otra salida
   - **"Ir a la lista de despachos"** — Para ver todas tus salidas

---

### Paso 8: Revisa tus despachos en la lista

*Un buen general revisa sus filas después de la batalla.*

1. Desde el menú lateral, haz clic en **"Registrar salidas"** o en el botón **"Ir a la lista de despachos"** tras guardar.

   *[Inserta aquí la captura de pantalla de la lista de despachos]*

2. Verás tus despachos organizados con tarjetas resumen. Cada tarjeta muestra:
   - **Número de salida** — Identificador único
   - **Tipo de movimiento** — Venta, devolución, etc.
   - **Fecha** — Cuándo se registró
   - **Estado** — Borrador, Confirmado o Cancelado
   - **Valor total** — Monto de la operación

3. Puedes **expandir** cualquier tarjeta para ver el detalle completo (productos, destinatario, bodega, observaciones).

4. Usa la barra de **búsqueda** para encontrar despachos por número, tipo o destinatario.

5. Usa el **filtro de estado** para ver solo Borradores, Confirmados o Cancelados.

   *[Inserta aquí la captura de pantalla de una tarjeta expandida con detalle]*

---

### Paso 9: Administra tus despachos

*Paulistas enseñaba que el honor está en mantener la palabra. Si algo cambia, actúa.*

- **Eliminar un despacho** — Si cometiste un error, puedes eliminarlo. Haz clic en el botón **"Eliminar"** dentro de la tarjeta expandida y confirma con **"¿Está seguro de eliminar este despacho?"**.

  | Mensaje sugerido |
  |------------------|
  | ⚠️ **"¿Está seguro de eliminar este despacho? Esta acción no se puede deshacer."** |

> **Importante:** Eliminar un despacho no recupera el inventario automáticamente en todos los casos. Si necesitas ajustar el stock, considera crear un ajuste de inventario.

---

## Lo que sigue: tu próxima misión

Con tu primera salida registrada, ya estás utilizando Paulu para gobernar tu negocio. Estos son tus siguientes pasos:

- ➡️ **Crea órdenes de compra** — Desde "Registrar orden" puedes solicitar productos a tus proveedores.
- ➡️ **Recibe mercancía** — Usa "Registrar entradas" para dar entrada a los productos que llegan.
- ➡️ **Aprovisiona tu tienda** — Crea transferencias entre bodegas para mover productos de tu almacén a la tienda.
- ➡️ **Monitorea tu inventario** — Revisa los productos y su stock mínimo para saber cuándo reordenar.

---

## Resolución de problemas (FAQ)

### No aparece mi producto en la lista al agregarlo
Asegúrate de haber creado el producto primero desde **"Registrar productos"** en el menú lateral. Si ya lo creaste, verifica que esté asociado a la empresa correcta.

### No encuentro el destinatario que busco
Regístralo primero desde **"Registrar proveedores"** (para proveedores) o verifica que el tipo de destinatario seleccionado sea el correcto.

### El botón "Registrar despacho" no hace nada
Revisa los campos marcados en rojo. El sistema requiere:
- Número de salida
- Fecha
- ID de bodega
- ID del destinatario
- Al menos un producto con cantidad mayor a 0

### Registré un despacho pero me equivoqué en la cantidad
Si el despacho está en estado **Borrador**, puedes eliminarlo y crear uno nuevo. Si está **Confirmado**, elimínalo y crea un ajuste de inventario para corregir el stock.

### ¿Puedo editar un despacho después de guardarlo?
Actualmente no. Si necesitas cambios, elimina el despacho y crea uno nuevo. Por eso recomendamos usar el estado **Borrador** hasta estar seguros.

### ¿Qué significa "Merma"?
Es el registro de productos que se dañaron, vencieron o se perdieron. Úsalo para mantener tu inventario preciso.

### ¿Cómo sé cuánto IVA debo declarar?
El sistema calcula el 19% de IVA automáticamente en cada despacho. Al final del mes, suma los totales de tus despachos tipo "Venta" para conocer tu base gravable.

---

## Mensajes de interfaz sugeridos (Microcopy)

### Botones y acciones principales

| Acción | Texto sugerido |
|--------|----------------|
| Registrar salida | **"Registrar despacho"** |
| Guardando... | **"Guardando..."** (con spinner) |
| Volver a la lista | **"Volver a la lista"** |
| Crear otro | **"Crear otro despacho"** |
| Ir a la lista | **"Ir a la lista de despachos"** |
| Agregar producto | **"+ Agregar producto"** |
| Generar código | **"Generar"** |
| Seleccionar destinatario | **"Seleccionar"** |
| Cancelar acción | **"Cancelar"** |
| Eliminar | **"Eliminar"** |
| Cerrar modal | **"Cerrar"** |

### Mensajes de éxito

| Contexto | Microcopy |
|----------|-----------|
| Despacho creado | ✅ **"¡Despacho registrado!"** |
| | *"Todo en orden. El despacho ha quedado registrado en el sistema."* |
| | *"Movimiento de salida registrado correctamente. El inventario ha sido actualizado."* |
| | *"Despacho creado. Los productos han sido asignados a su destino."* |
| Operación general exitosa | ✅ **"Operación completada con éxito."** |

### Mensajes de advertencia y error

| Contexto | Microcopy |
|----------|-----------|
| Error al guardar | ⚠️ **"Error al registrar el despacho."** *Revisa los datos e intenta de nuevo.* |
| Campo obligatorio faltante | ⚠️ **"El número de despacho es obligatorio."** |
| | ⚠️ **"La fecha es obligatoria."** |
| | ⚠️ **"El ID de bodega es obligatorio."** |
| | ⚠️ **"El ID del destinatario es obligatorio."** |
| Sin productos | ⚠️ **"Debe agregar al menos un producto."** |
| Cantidad inválida | ⚠️ **"La cantidad debe ser mayor a 0."** |
| Costo inválido | ⚠️ **"El costo unitario no puede ser negativo."** |
| Confirmación de eliminación | ⚠️ **"¿Está seguro de eliminar este despacho?"** |
| Sin resultados en búsqueda | **"No se encontraron despachos."** |
| Sin productos en inventario | **"Aún no tienes productos registrados. Cada producto que añadas hoy es una venta que podrás controlar mañana."** |
| Error de conexión | ⚠️ **"No pudimos conectar con el servidor. Verifica tu conexión e intenta de nuevo."** |

### Textos de ayuda en formularios

| Campo | Texto de ayuda |
|-------|----------------|
| Número de salida | "Formato: SAL-20260615-0001" |
| Cantidad del producto | "Ingresa solo números mayores a 0" |
| Costo unitario | "Precio de venta por unidad" |
| Observaciones | "Notas adicionales sobre esta salida (opcional)" |
| Búsqueda en lista | "Buscar por número, tipo o destinatario..." |

---

> *"El valor de un hombre se mide por su honor, su lealtad y su responsabilidad hacia quienes lo siguen."* — Paulus Atreides

Cada despacho que registras, cada producto que controlas y cada venta que concretas es un acto de liderazgo sobre tu negocio. Paulu está aquí para que gobiernes tu inventario con la misma determinación con la que el Duque Paulus comandaba sus tropas: desde el frente, sin dudar y con honor.

**Bienvenido al frente de batalla. Bienvenido a Paulu.**

---

*Documento v1.0 — Junio 2026*
